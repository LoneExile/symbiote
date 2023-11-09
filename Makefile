
build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o symbiote main.go
	rm -rf ~/go/bin/symbiote
	mv symbiote ~/go/bin/symbiote
	mkdir -p ~/.config/symbiote
	cp config.yaml ~/.config/symbiote/config.yaml
	~/go/bin/symbiote completion zsh > ~/.config/symbiote/symbiote.zsh
	# echo "source ~/.config/symbiote/symbiote.zsh" >> ~/.zshrc
	# ~/go/bin/symbiote completion bash > ~/.config/symbiote/symbiote.bash
	# echo "source ~/.config/symbiote/symbiote.bash" >> ~/.bashrc
