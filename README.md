![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/CGI-FR/MIMO/ci.yml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cgi-fr/mimo)](https://goreportcard.com/report/github.com/cgi-fr/mimo)
![GitHub all releases](https://img.shields.io/github/downloads/CGI-FR/MIMO/total)
![GitHub](https://img.shields.io/github/license/CGI-FR/MIMO)
![GitHub Repo stars](https://img.shields.io/github/stars/CGI-FR/MIMO)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/CGI-FR/MIMO)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CGI-FR/MIMO)

# MIMO : Masked Input Metrics Output

## Usage

```console
> mkfifo real.jsonl # create a pipe file to store the real json stream before pseudonymization
> lino pull prod | tee real.jsonl | pimo | mimo real.jsonl | lino push dev
8:27AM WRN field is not completely masked fieldname=surname

       MIMO REPORT
===========================================
fieldname | masking rate | collision rate |
----------|--------------|----------------|
name      |        100 % |            0 % |
surname   |         99 % |            0 % |
> rm real.jsonl # pipe file can be removed after
```

```bash
pimo --empty-input --repeat 1000 --mask 'name=[{add:""},{randomChoiceInUri:"pimo://nameFR"}]' | tee real.jsonl | pimo --mask 'name={randomChoiceInUri:"pimo://nameFR"}' | mimo real.jsonl
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
