package stencil_app

import (
	"fmt"
	"github.com/arelate/gaugin/data"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"strconv"
	"strings"
)

const (
	transitiveOpen  = " ("
	transitiveClose = ")"
)

var labelTitles = map[string]string{
	vangogh_local_data.OwnedProperty:          "Own",
	vangogh_local_data.TBAProperty:            "TBA",
	vangogh_local_data.ComingSoonProperty:     "Soon",
	vangogh_local_data.PreOrderProperty:       "PO",
	vangogh_local_data.InDevelopmentProperty:  "In Dev",
	vangogh_local_data.IsUsingDOSBoxProperty:  "DOSBox",
	vangogh_local_data.IsUsingScummVMProperty: "ScummVM",
	vangogh_local_data.IsFreeProperty:         "Free",
	vangogh_local_data.WishlistedProperty:     "Wish",
}

func transitiveDst(s string) string {
	dst := s
	if strings.Contains(s, transitiveOpen) {
		dst = s[:strings.LastIndex(s, transitiveOpen)]
	}
	return dst
}

func transitiveSrc(s string) string {
	src := ""
	if strings.Contains(s, transitiveOpen) {
		from, to := strings.LastIndex(s, transitiveOpen)+len(transitiveOpen), strings.Index(s, transitiveClose)
		if from > to {
			to = strings.LastIndex(s, transitiveClose)
			if from > to {
				from = 0
				to = len(s) - 1
			}
		}
		src = s[from:to]
	}
	return src
}

func ownedValidationResult(id string, rxa kvas.ReduxAssets) (string, bool) {
	return rxa.GetFirstVal(vangogh_local_data.ValidationResultProperty, id)
}

func ReviewClass(sr string) string {
	if strings.Contains(sr, "Positive") {
		return "positive"
	} else if strings.Contains(sr, "Negative") {
		return "negative"
	} else {
		return "neutral"
	}
}

func fmtAction(id, property, link string, rxa kvas.ReduxAssets) string {

	owned, _ := rxa.GetFirstVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		if owned == "true" {
			return ""
		}
		switch link {
		case "true":
			return "Remove"
		case "false":
			return "Add"
		}
	case vangogh_local_data.TagIdProperty:
		if owned != "true" {
			return ""
		}
		return "Edit"
	case vangogh_local_data.LocalTagsProperty:
		return "Edit"
	}
	return ""
}

func fmtActionHref(id, property, link string, rxa kvas.ReduxAssets) string {

	owned, _ := rxa.GetFirstVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		switch link {
		case "Add":
			return "/wishlist/add?id=" + id
		case "Remove":
			return "/wishlist/remove?id=" + id
		}
	case vangogh_local_data.TagIdProperty:
		if owned != "true" {
			return ""
		}
		return "/tags/edit?id=" + id
	case vangogh_local_data.LocalTagsProperty:
		return "/local-tags/edit?id=" + id
	}
	return ""
}

func fmtClass(id, property, link string, rxa kvas.ReduxAssets) string {
	switch property {
	case vangogh_local_data.OwnedProperty:
		if res, ok := ownedValidationResult(id, rxa); ok {
			if res == "OK" {
				return "validation-result-ok"
			} else {
				return "validation-result-err"
			}
		} else {
			return ""
		}
	case vangogh_local_data.SteamReviewScoreDescProperty:
		return ReviewClass(link)
	case vangogh_local_data.RatingProperty:
		return ReviewClass(fmtGOGRating(link))
	}
	return ""
}

func fmtHref(_, property, link string, _ kvas.ReduxAssets) string {
	switch property {
	case vangogh_local_data.GOGOrderDateProperty:
		link = justTheDate(link)
	case vangogh_local_data.PublishersProperty:
		fallthrough
	case vangogh_local_data.DevelopersProperty:
		return fmt.Sprintf("/search?%s=%s&sort=global-release-date&desc=true", property, link)
	case vangogh_local_data.IncludesGamesProperty:
		fallthrough
	case vangogh_local_data.IsIncludedByGamesProperty:
		fallthrough
	case vangogh_local_data.RequiresGamesProperty:
		fallthrough
	case vangogh_local_data.IsRequiredByGamesProperty:
		return fmt.Sprintf("/product?id=%s", transitiveSrc(link))
	case vangogh_local_data.RatingProperty:
		return ""
	case vangogh_local_data.DiscountPercentageProperty:
		return ""
	case vangogh_local_data.PriceProperty:
		return ""
	case data.GauginGOGLinksProperty:
		fallthrough
	case data.GauginSteamLinksProperty:
		return transitiveSrc(link)
	}
	return fmt.Sprintf("/search?%s=%s", property, link)
}

func justTheDate(s string) string {
	return strings.Split(s, " ")[0]
}

func fmtLabel(id, property, link string, rxa kvas.ReduxAssets) string {

	label := link
	owned, _ := rxa.GetFirstVal(vangogh_local_data.OwnedProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		fallthrough
	case vangogh_local_data.OwnedProperty:
		fallthrough
	case vangogh_local_data.PreOrderProperty:
		fallthrough
	case vangogh_local_data.ComingSoonProperty:
		fallthrough
	case vangogh_local_data.TBAProperty:
		fallthrough
	case vangogh_local_data.InDevelopmentProperty:
		fallthrough
	case vangogh_local_data.IsUsingDOSBoxProperty:
		fallthrough
	case vangogh_local_data.IsUsingScummVMProperty:
		fallthrough
	case vangogh_local_data.IsFreeProperty:
		if link == "true" {
			return labelTitles[property]
		}
		return ""
	case vangogh_local_data.ProductTypeProperty:
		if link == "GAME" {
			return ""
		}
	case vangogh_local_data.DiscountPercentageProperty:
		if owned == "true" {
			return ""
		}
		if link != "" && link != "0" {
			return fmt.Sprintf("-%s%%", link)
		}
		return ""
	case vangogh_local_data.TagIdProperty:
		return transitiveDst(link)
	}
	return label
}

func fmtTitle(id, property, link string, rxa kvas.ReduxAssets) string {
	title := link

	owned, _ := rxa.GetFirstVal(vangogh_local_data.OwnedProperty, id)
	isFree, _ := rxa.GetFirstVal(vangogh_local_data.IsFreeProperty, id)

	switch property {
	case vangogh_local_data.WishlistedProperty:
		if owned == "true" {
			return ""
		}
		return DigestTitles[link]
	case vangogh_local_data.IncludesGamesProperty:
		fallthrough
	case vangogh_local_data.IsIncludedByGamesProperty:
		fallthrough
	case vangogh_local_data.RequiresGamesProperty:
		fallthrough
	case vangogh_local_data.IsRequiredByGamesProperty:
		title = transitiveDst(link)
	case vangogh_local_data.GOGOrderDateProperty:
		title = justTheDate(link)
	case vangogh_local_data.LanguageCodeProperty:
		title = LanguageCodeFlag(transitiveSrc(link)) + " " + transitiveDst(link)
	case vangogh_local_data.RatingProperty:
		title = fmtGOGRating(link)
	case vangogh_local_data.TagIdProperty:
		return transitiveDst(link)
	case vangogh_local_data.PriceProperty:
		if isFree == "true" {
			return ""
		}
	case data.GauginGOGLinksProperty:
		fallthrough
	case data.GauginSteamLinksProperty:
		title = PropertyTitles[transitiveDst(link)]
	}

	return title
}

func fmtGOGRating(rs string) string {
	rd := ""
	if ri, err := strconv.ParseInt(rs, 10, 32); err == nil {
		if ri >= 45 {
			rd = "Very Positive"
		} else if ri > 35 {
			rd = "Positive"
		} else if ri > 25 {
			rd = "Mixed"
		} else if ri > 15 {
			rd = "Negative"
		} else if ri > 0 {
			rd = "Very Negative"
		} else {
			rd = "Not rated"
		}
		if ri > 0 {
			rd += fmt.Sprintf(" (%.1f)", float32(ri)/10.0)
		}
	}
	return rd
}