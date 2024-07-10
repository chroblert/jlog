package jlog

// GetAllWritedSize 获取本次在所有文件中已经写入的大小
func GetAllWritedSize() int64 {
	fishLogger.mu.Lock()
	defer fishLogger.mu.Unlock()
	return fishLogger.writed_size
}

// GetCurrentFileSize 获取在当前文件中已经写入的大小
func GetCurrentFileSize() int64 {
	fishLogger.mu.Lock()
	defer fishLogger.mu.Unlock()
	return fishLogger.size
}
