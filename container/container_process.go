package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"mydocker/constant"

	log "github.com/sirupsen/logrus"
)

const (
	RUNNING       = "running"
	STOP          = "stopped"
	Exit          = "exited"
	InfoLoc       = "/var/run/mydocker/"
	InfoLocFormat = InfoLoc + "%s/"
	ConfigName    = "config.json"
	IDLength      = 10
	LogFile       = "container.log"
)

type Info struct {
	Pid         string `json:"pid"`        // 容器的init进程在宿主机上的 PID
	Id          string `json:"id"`         // 容器Id
	Name        string `json:"name"`       // 容器名
	Command     string `json:"command"`    // 容器内init运行命令
	CreatedTime string `json:"createTime"` // 创建时间
	Status      string `json:"status"`     // 容器的状态
}

// NewParentProcess 构建 command 用于启动一个新进程
/*
这里是父进程，也就是当前进程执行的内容。
1.这里的/proc/se1f/exe调用中，/proc/self/ 指的是当前运行进程自己的环境，exec 其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化
2.后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用initCommand去初始化进程的一些环境和资源
3.下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境。
4.如果用户指定了-it参数，就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, volume, containerName string) (*exec.Cmd, *os.File) {
	// 创建匿名管道用于传递参数，将readPipe作为子进程的ExtraFiles，子进程从readPipe中读取参数
	// 父进程中则通过writePipe将参数写入管道
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// 对于后台运行容器，将标准输出重定向到日志文件中，便于后续查询
		dirURL := fmt.Sprintf(InfoLocFormat, containerName)
		if err := os.MkdirAll(dirURL, constant.Perm0622); err != nil {
			log.Errorf("NewParentProcess mkdir %s error %v", dirURL, err)
			return nil, nil
		}
		stdLogFilePath := dirURL + LogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	mntURL := "/root/merged"
	rootURL := "/root"
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL
	return cmd, writePipe
}

// NewWorkSpace Create an Overlay2 filesystem as container root workspace
func NewWorkSpace(rootURL, mntURL, volume string) {
	log.Infof("createLower")
	createLower(rootURL)
	createDirs(rootURL)
	mountOverlayFS(rootURL, mntURL)
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			log.Infof("mountVolume")
			mountVolume(mntURL, volumeURLs)
			log.Infof("volumeURL:%v", volumeURLs)
		} else {
			log.Infof("volume parameter input is not correct.")
		}
	}
}

func mountVolume(mntURL string, volumeURLs []string) {
	// 第0个元素为宿主机目录
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, constant.Perm0777); err != nil {
		log.Infof("mkdir parent dir %s error. %v", parentUrl, err)
	}
	// 第1个元素为容器目录
	containerUrl := volumeURLs[1]
	// 拼接并创建对应的容器目录
	containerVolumeURL := mntURL + "/" + containerUrl
	if err := os.Mkdir(containerVolumeURL, constant.Perm0777); err != nil {
		log.Infof("mkdir container dir %s error. %v", containerVolumeURL, err)
	}
	// 通过bind mount 将宿主机目录挂载到容器目录
	// mount -o bind /hostURL /containerURL
	cmd := exec.Command("mount", "-o", "bind", parentUrl, containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("mount volume failed. %v", err)
	}
}

// createLower 将busybox作为overlayfs的lower层
func createLower(rootURL string) {
	// 把busybox作为overlayfs中的lower层
	busyboxURL := rootURL + "/busybox"
	busyboxTarURL := rootURL + "/busybox.tar"
	// 检查是否已经存在busybox文件夹
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	// 不存在则创建目录并将busybox.tar解压到busybox文件夹中
	if !exist {
		if err := os.Mkdir(busyboxURL, constant.Perm0777); err != nil {
			log.Errorf("mkdir dir %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("untar dir %s error %v", busyboxURL, err)
		}
	}
}

// createDirs 创建overlayfs需要的的upper、worker目录
func createDirs(rootURL string) {
	upperURL := rootURL + "/upper"
	if err := os.Mkdir(upperURL, constant.Perm0777); err != nil {
		log.Errorf("mkdir dir %s error. %v", upperURL, err)
	}
	workURL := rootURL + "/work"
	if err := os.Mkdir(workURL, constant.Perm0777); err != nil {
		log.Errorf("mkdir dir %s error. %v", workURL, err)
	}
}

// mountOverlayFS 挂载overlayfs
func mountOverlayFS(rootURL string, mntURL string) {
	// 创建对应的挂载目录
	if err := os.Mkdir(mntURL, constant.Perm0777); err != nil {
		log.Errorf("mountOverlayFS mkdir dir %s error. %v", mntURL, err)
	}
	// 拼接参数
	// e.g. lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/merged
	dirs := "lowerdir=" + rootURL + "/busybox" + ",upperdir=" + rootURL + "/upper" + ",workdir=" + rootURL + "/work"

	// 完整命令：mount -t overlay overlay -o lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/work /root/merged
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mntURL)
	log.Infof("mountOverlayFS cmd:%s", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("mountOverlayFS mount err:%v", err)
	}
}

// DeleteWorkSpace Delete the OverlayFS filesystem while container exit
func DeleteWorkSpace(rootURL, mntURL, volume string) {
	// 如果指定了volume则需要先umount volume
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			umountVolume(mntURL, volumeURLs)
		}
	}
	// 然后umount整个容器的挂载点
	umountOverlayFS(mntURL)
	// 最后移除相关文件夹
	deleteDirs(rootURL)
}

func umountOverlayFS(mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("umountOverlayFS umount %s err:%v", mntURL, err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove dir %s error %v", mntURL, err)
	}
}

func umountVolume(mntURL string, volumeURLs []string) {
	containerUrl := mntURL + "/" + volumeURLs[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount volume failed. %v", err)
	}
}

func deleteDirs(rootURL string) {
	writeURL := rootURL + "/upper"
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("remove dir %s error %v", writeURL, err)
	}
	workURL := rootURL + "/work"
	if err := os.RemoveAll(workURL); err != nil {
		log.Errorf("remove dir %s error %v", workURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// volumeUrlExtract 通过冒号分割解析volume目录，比如 -v /tmp:/tmp
func volumeUrlExtract(volume string) []string {
	return strings.Split(volume, ":")
}
