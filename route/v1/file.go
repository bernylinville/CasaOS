package v1

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/file"
	oasis_err2 "github.com/IceWhaleTech/CasaOS/pkg/utils/oasis_err"
	"github.com/IceWhaleTech/CasaOS/service"
	"github.com/gin-gonic/gin"
)

func downloadReadFile(c *gin.Context) {
	//http下载地址 csv
	csvFileUrl := c.PostForm("file_name")
	res, err := http.Get(csvFileUrl)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	defer res.Body.Close()
	//读取csv
	reader := csv.NewReader(bufio.NewReader(res.Body))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			c.String(400, err.Error())
			return
		}
		//line 就是每一行的内容
		fmt.Println(line)
		//line[0] 就是第几列
		fmt.Println(line[0])
	}
}

func downloadWriteFile(c *gin.Context) {
	//写文件
	var filename = "./output1.csv"

	file, err := os.Create(filename) //创建文件
	if err != nil {
		c.String(400, err.Error())
		return
	}
	buf := bufio.NewWriter(file) //创建新的 Writer 对象
	buf.WriteString("test")
	buf.Flush()
	defer file.Close()

	//返回文件流
	c.File(filename)
}

// @Summary 读取文件
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string true "路径"
// @Success 200 {string} string "ok"
// @Router /file/read [get]
func GetFilerContent(c *gin.Context) {
	filePath := c.Query("path")
	if len(filePath) == 0 {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.INVALID_PARAMS,
			Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(filePath) {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_DOES_NOT_EXIST,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_DOES_NOT_EXIST),
		})
		return
	}
	//文件读取任务是将文件内容读取到内存中。
	info, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_READ_ERROR,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_READ_ERROR),
			Data:    err.Error(),
		})
		return
	}
	result := string(info)

	//返回结果
	c.JSON(http.StatusOK, model.Result{
		Success: oasis_err2.SUCCESS,
		Message: oasis_err2.GetMsg(oasis_err2.SUCCESS),
		Data:    result,
	})
}

func GetLocalFile(c *gin.Context) {
	path := c.Query("path")
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.INVALID_PARAMS,
			Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(path) {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_DOES_NOT_EXIST,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_DOES_NOT_EXIST),
		})
		return
	}
	c.File(path)
	return
}

// @Summary 下载文件
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string true "路径"
// @Success 200 {string} string "ok"
// @Router /file/download [get]
func GetDownloadFile(c *gin.Context) {
	filePath := c.Query("path")
	if len(filePath) == 0 {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.INVALID_PARAMS,
			Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS),
		})
		return
	}
	if !file.Exists(filePath) {
		c.JSON(http.StatusOK, model.Result{
			Success: oasis_err2.FILE_DOES_NOT_EXIST,
			Message: oasis_err2.GetMsg(oasis_err2.FILE_DOES_NOT_EXIST),
		})
		return
	}
	//打开文件
	fileTmp, _ := os.Open(filePath)
	defer fileTmp.Close()
	//获取文件的名称
	fileName := path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	c.File(filePath)
	return
}

// @Summary 获取目录列表
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param path query string false "路径"
// @Success 200 {string} string "ok"
// @Router /file/dirpath [get]
func DirPath(c *gin.Context) {
	path := c.DefaultQuery("path", "/")
	info := service.MyService.ZiMa().GetDirPath(path)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS), Data: info})
}

// @Summary 重命名目录或文件
// @Produce  application/json
// @Accept application/json
// @Tags file
// @Security ApiKeyAuth
// @Param oldpath formData string true "旧的路径"
// @Param newpath formData string true "新路径"
// @Success 200 {string} string "ok"
// @Router /file/rename [put]
func RenamePath(c *gin.Context) {
	op := c.PostForm("oldpath")
	np := c.PostForm("newpath")
	if len(op) == 0 || len(np) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	service.MyService.ZiMa().RenameFile(op, np)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

// @Summary 创建文件夹
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string false "路径"
// @Success 200 {string} string "ok"
// @Router /file/mkdir [post]
func MkdirAll(c *gin.Context) {
	path := c.PostForm("path")
	var code int
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	code, _ = service.MyService.ZiMa().MkdirAll(path)
	c.JSON(http.StatusOK, model.Result{Success: code, Message: oasis_err2.GetMsg(code)})
}

// @Summary 创建文件
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string false "路径"
// @Success 200 {string} string "ok"
// @Router /file/create [post]
func PostCreateFile(c *gin.Context) {
	path := c.PostForm("path")
	var code int
	if len(path) == 0 {
		c.JSON(http.StatusOK, model.Result{Success: oasis_err2.INVALID_PARAMS, Message: oasis_err2.GetMsg(oasis_err2.INVALID_PARAMS)})
		return
	}
	code, _ = service.MyService.ZiMa().CreateFile(path)
	c.JSON(http.StatusOK, model.Result{Success: code, Message: oasis_err2.GetMsg(code)})
}

// @Summary 上传文件
// @Produce  application/json
// @Accept  multipart/form-data
// @Tags file
// @Security ApiKeyAuth
// @Param path formData string false "路径"
// @Success 200 {string} string "ok"
// @Router /file/mkdir [post]
func PostFileUpload(c *gin.Context) {
	file, _, _ := c.Request.FormFile("file")
	//file.Read()
	path := c.Query("path")
	//上传文件
	out, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	defer out.Close()
	io.Copy(out, file)
	c.JSON(http.StatusOK, model.Result{Success: oasis_err2.SUCCESS, Message: oasis_err2.GetMsg(oasis_err2.SUCCESS)})
}

func PutFileMove(c *gin.Context) {
	from := "/Users/liangjianli/go/CasaOS"
	to := "/Users/liangjianli/go/CasaOS/test"
	//t := 1 //是否覆盖

	//方法体
	stopCh := make(chan int)
	f, err := os.Stat(from)
	if err != nil {
		//未拿到文件信息
		fmt.Println("stat", err)
	}
	//未创建新的文件夹
	if f.IsDir() {
		//from 是文件夹,定义to也是文件夹
		if list, err := ioutil.ReadDir(from); err == nil {
			for _, v := range list {
				time.Sleep(time.Second)
				if err = Copy(stopCh, filepath.Join(from, v.Name()), filepath.Join(to, v.Name())); err != nil {
					fmt.Printf("copy %s ,err %d", v.Name(), err)
				}
			}
		}
	} else {
		p := filepath.Dir(to)
		if _, err = os.Stat(p); err != nil {
			if err = os.MkdirAll(p, 0777); err != nil {
				fmt.Println("mkdir", err)
			}
		}
	}

	file, err := os.Open(from)

	if err != nil {
		fmt.Println("open file error ", err)
	}
	defer file.Close()
	out, err := os.Create(to)
	if err != nil {
		fmt.Println("create to file err", err)
	}
	defer out.Close()
	io.Copy(out, file)
	time.Sleep(time.Second * 4)
	close(stopCh)
}
func Copy(stop chan int, from, to string) error {

	for {
		select {
		case <-stop:
			return nil
		default:
			fmt.Println(from)

		}
	}
	return nil
}
