# MIMO

Sensitive and
Personal
Information
Protection
Evaluator

Secure Pseudo-Anonymization Metrics

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

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
