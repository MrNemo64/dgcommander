package dgc

import (
	"reflect"

	"github.com/bwmarrin/discordgo"
)

func removeElement[T comparable](slice []T, val T) []T {
	for i, v := range slice {
		if v == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func findOption(name string, options []*discordgo.ApplicationCommandInteractionDataOption) *discordgo.ApplicationCommandInteractionDataOption {
	for _, v := range options {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func missingKeys[K comparable, V any](reference []K, actual map[K]V) []K {
	var missing []K
	for i := range reference {
		if _, found := actual[reference[i]]; !found {
			missing = append(missing, reference[i])
		}
	}
	return missing
}

func nameOfT[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}
