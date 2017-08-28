# The impact of inheritance on open-ended problem solving

This project investigates the impact of inheritance on problem solving ability in the scientific discovery game **FoldIt**. Problem solving in **FoldIt** involves folding amino acid structures into three-dimensional protein structures. Individuals may work alone or together by sharing partial solutions for others to inherit. This investigation measures the impact of inheritance on problem solving ability by comparing the effectiveness of teams of problem solvers to individuals working alone. The goal of this investigation is to uncover the tradeoffs involved in sharing partial solutions to problems among groups of problem solvers in open-ended problem solving contexts.

# Repo

This repo contains go, python, and R code.

```bash
go get github.com/pedmiston/foldit
pip install git+git://github.com/pedmiston/foldit#egg=foldit
Rscript -e "devtools::install_github('pedmiston/foldit')"
```

# FoldIt server

Requirements:

- awscli

## AWS

To configure AWS for S3 storage, install the `awscli` and run `aws configure`.
