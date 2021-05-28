package validationtranslatormock

//go:generate mockgen -destination=mock.go -package=$GOPACKAGE github.com/KennyChenFight/Shortening-URL/pkg/validation Translator
