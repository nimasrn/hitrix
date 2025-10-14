package cdn

import (
	"strings"

	"github.com/coretrix/hitrix/service"
)

const (
	Options = "fit=${fit},format=${format},metadata=none,onerror=redirect,quality=${quality},width=${width},dpr=${dpr}/"
)

func GetImageURLTemplate(image string) string {
	return service.DI().Config().MustString("oss.cdn_url") + Options + image
}

func GetImageURLTemplateFilled(image, fit, format, quality, width, dpr string) string {
	image = service.DI().Config().MustString("oss.cdn_url") + Options + image
	image = strings.ReplaceAll(image, "${fit}", fit)
	image = strings.ReplaceAll(image, "${format}", format)
	image = strings.ReplaceAll(image, "${quality}", quality)
	image = strings.ReplaceAll(image, "${width}", width)
	image = strings.ReplaceAll(image, "${dpr}", dpr)

	return image
}
