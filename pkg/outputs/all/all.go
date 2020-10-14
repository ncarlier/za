package all

import (
	// activate file output writer
	_ "github.com/ncarlier/za/pkg/outputs/file"
	// activate Prometheus output writer
	// _ "github.com/ncarlier/za/pkg/outputs/prometheus"
	// activate Loki output writer
	_ "github.com/ncarlier/za/pkg/outputs/loki"
)
