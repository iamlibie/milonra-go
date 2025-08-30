package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iamlibie/milonra-go/plugin"
)

// UploadGroupFile 上传群文件
func UploadGroupFile(b plugin.Bot, groupID int64, file, name string, folder ...string) error {
	echo := generateEcho("upload_group_file")
	params := map[string]interface{}{
		"group_id": groupID,
		"file":     file,
		"name":     name,
	}

	if len(folder) > 0 {
		params["folder"] = folder[0]
	} else {
		params["folder"] = "/"
	}

	data := map[string]interface{}{
		"action": "upload_group_file",
		"params": params,
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("上传群文件失败: %s", resp.Status)
	}

	return nil
}

// UploadPrivateFile 上传私聊文件
func UploadPrivateFile(b plugin.Bot, userID int64, file, name string) error {
	echo := generateEcho("upload_private_file")
	data := map[string]interface{}{
		"action": "upload_private_file",
		"params": map[string]interface{}{
			"user_id": userID,
			"file":    file,
			"name":    name,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("上传私聊文件失败: %s", resp.Status)
	}

	return nil
}

// GetGroupFileURL 获取群文件资源链接
func GetGroupFileURL(b plugin.Bot, groupID int64, fileID string, busid ...int) (string, error) {
	echo := generateEcho("get_group_file_url")
	params := map[string]interface{}{
		"group_id": groupID,
		"file_id":  fileID,
	}

	if len(busid) > 0 {
		params["busid"] = busid[0]
	}

	data := map[string]interface{}{
		"action": "get_group_file_url",
		"params": params,
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return "", err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return "", err
	}

	if resp.Status != "ok" {
		return "", fmt.Errorf("获取群文件链接失败: %s", resp.Status)
	}

	var result struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

// GetPrivateFileURL 获取私聊文件资源链接
func GetPrivateFileURL(b plugin.Bot, userID int64, fileID string, fileHash ...string) (string, error) {
	echo := generateEcho("get_private_file_url")
	params := map[string]interface{}{
		"user_id": userID,
		"file_id": fileID,
	}

	if len(fileHash) > 0 {
		params["file_hash"] = fileHash[0]
	}

	data := map[string]interface{}{
		"action": "get_private_file_url",
		"params": params,
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return "", err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return "", err
	}

	if resp.Status != "ok" {
		return "", fmt.Errorf("获取私聊文件链接失败: %s", resp.Status)
	}

	var result struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

// GroupFileInfo 群文件信息
type GroupFileInfo struct {
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	BusID        int    `json:"busid"`
	FileSize     int64  `json:"file_size"`
	UploadTime   int64  `json:"upload_time"`
	DeadTime     int64  `json:"dead_time"`
	ModifyTime   int64  `json:"modify_time"`
	DownloadTime int    `json:"download_times"`
	Uploader     int64  `json:"uploader"`
	UploaderName string `json:"uploader_name"`
}

// GroupFolderInfo 群文件夹信息
type GroupFolderInfo struct {
	FolderID       string `json:"folder_id"`
	FolderName     string `json:"folder_name"`
	CreateTime     int64  `json:"create_time"`
	Creator        int64  `json:"creator"`
	CreatorName    string `json:"creator_name"`
	TotalFileCount int    `json:"total_file_count"`
}

// GetGroupRootFiles 获取群根目录文件列表
func GetGroupRootFiles(b plugin.Bot, groupID int64) (*GroupFilesData, error) {
	echo := generateEcho("get_group_root_files")
	data := map[string]interface{}{
		"action": "get_group_root_files",
		"params": map[string]interface{}{
			"group_id": groupID,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return nil, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("获取群文件列表失败: %s", resp.Status)
	}

	var filesData GroupFilesData
	err = json.Unmarshal(resp.Data, &filesData)
	if err != nil {
		return nil, err
	}

	return &filesData, nil
}

// GroupFilesData 群文件列表数据
type GroupFilesData struct {
	Files   []GroupFileInfo   `json:"files"`
	Folders []GroupFolderInfo `json:"folders"`
}

// GetGroupFilesByFolder 获取群子目录文件列表
func GetGroupFilesByFolder(b plugin.Bot, groupID int64, folderID string) (*GroupFilesData, error) {
	echo := generateEcho("get_group_files_by_folder")
	data := map[string]interface{}{
		"action": "get_group_files_by_folder",
		"params": map[string]interface{}{
			"group_id":  groupID,
			"folder_id": folderID,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return nil, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("获取群文件夹内容失败: %s", resp.Status)
	}

	var filesData GroupFilesData
	err = json.Unmarshal(resp.Data, &filesData)
	if err != nil {
		return nil, err
	}

	return &filesData, nil
}

// CreateGroupFileFolder 创建群文件文件夹（只能在根目录创建）
func CreateGroupFileFolder(b plugin.Bot, groupID int64, name string) error {
	echo := generateEcho("create_group_file_folder")
	data := map[string]interface{}{
		"action": "create_group_file_folder",
		"params": map[string]interface{}{
			"group_id":  groupID,
			"name":      name,
			"parent_id": "/", // TX不再允许在非根目录创建文件夹
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("创建文件夹失败: %s", resp.Status)
	}

	return nil
}

// DeleteGroupFile 删除群文件
func DeleteGroupFile(b plugin.Bot, groupID int64, fileID string) error {
	echo := generateEcho("delete_group_file")
	data := map[string]interface{}{
		"action": "delete_group_file",
		"params": map[string]interface{}{
			"group_id": groupID,
			"file_id":  fileID,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("删除文件失败: %s", resp.Status)
	}

	return nil
}

// DeleteGroupFileFolder 删除群文件文件夹
func DeleteGroupFileFolder(b plugin.Bot, groupID int64, folderID string) error {
	echo := generateEcho("delete_group_file_folder")
	data := map[string]interface{}{
		"action": "delete_group_file_folder",
		"params": map[string]interface{}{
			"group_id":  groupID,
			"folder_id": folderID,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("删除文件夹失败: %s", resp.Status)
	}

	return nil
}

// MoveGroupFile 移动群文件
func MoveGroupFile(b plugin.Bot, groupID int64, fileID, parentDir, targetDir string) error {
	echo := generateEcho("move_group_file")
	data := map[string]interface{}{
		"action": "move_group_file",
		"params": map[string]interface{}{
			"group_id":         groupID,
			"file_id":          fileID,
			"parent_directory": parentDir,
			"target_directory": targetDir,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("移动文件失败: %s", resp.Status)
	}

	return nil
}

// RenameGroupFileFolder 重命名群文件文件夹
func RenameGroupFileFolder(b plugin.Bot, groupID int64, folderID, newName string) error {
	echo := generateEcho("rename_group_file_folder")
	data := map[string]interface{}{
		"action": "rename_group_file_folder",
		"params": map[string]interface{}{
			"group_id":        groupID,
			"folder_id":       folderID,
			"new_folder_name": newName,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("重命名文件夹失败: %s", resp.Status)
	}

	return nil
}

// UploadImage 上传图片
func UploadImage(b plugin.Bot, file string) (string, error) {
	echo := generateEcho("upload_image")
	data := map[string]interface{}{
		"action": "upload_image",
		"params": map[string]interface{}{
			"file": file,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return "", err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return "", err
	}

	if resp.Status != "ok" {
		return "", fmt.Errorf("上传图片失败: %s", resp.Status)
	}

	// 返回的data直接是文件URL字符串
	var url string
	err = json.Unmarshal(resp.Data, &url)
	if err != nil {
		return "", err
	}

	return url, nil
}

// 检查文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
