package domain

type GrafanaAlert struct {
	Reciver           string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []Alert           `json:"alerts"`
	CommonLabels      CommonLabels      `json:"commonLabels"`
	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
	Title             string            `json:"title"`
}

type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Anotations  map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt"`
	EndsAt      string            `json:"endsAt"`
	DashboarUrl string            `json:"dashboardUrl"`
	SilenceURL  string            `json:"silenceUrl"`
	Values      map[string]any    `json:"values"`
	ValueString string            `json:"valueString"`
}

type CommonLabels struct {
	Alertname     string `json:"alertname"`
	GrafanaFolder string `json:"grafana_folder"`
}

type CommonAnnotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}
