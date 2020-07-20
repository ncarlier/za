package all

import (
	// activate file output writer
	_ "github.com/ncarlier/trackr/pkg/outputs/file"
	// activate Prometheus output writer
	_ "github.com/ncarlier/trackr/pkg/outputs/prometheus"
	// activate Loki output writer
	_ "github.com/ncarlier/trackr/pkg/outputs/loki"
)
