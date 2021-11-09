# searchHIBP

searchHIBP is a golang tool that implements binary search over a hash ordered binary file.
You may create one with https://github.com/fblz/convertHIBP
Enter your password on the commandline and it will tell you whether it is included in the InputFile.

## Installation

```
git clone https://github.com/fblz/searchHIBP.git
cd searchHIBP
go build
#  <or>
go install
```

## Usage

```
searchHIBP
  -InputFile string
        Specify the binary hibp file (default "none")
```