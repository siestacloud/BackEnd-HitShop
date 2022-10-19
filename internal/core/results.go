package core

// ExtractResult итоговый результат по извлечению сессий из новых аккаунтов
type ExtractResult struct {
	TotalFiles             int `json:"total_files"`
	TotalExtractedSessions int `json:"total_extracted_sessions"`
	TotalValidSessions     int `json:"total_valid_sessions"`
}
