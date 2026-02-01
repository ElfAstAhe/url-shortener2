package dto

type InternalStatsDto struct {
	TotalURLCount  int `json:"urls"`
	TotalUserCount int `json:"users"`
}

func NewInternalStatsDto(totalURLCount int, userCount int) *InternalStatsDto {
	return &InternalStatsDto{
		TotalURLCount:  totalURLCount,
		TotalUserCount: userCount,
	}
}
