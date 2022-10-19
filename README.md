# wiki_demo

### Demo Wiki API Tool


#### Clone
```bash
git clone https://github.com/solutionstack/wiki_demo.git
```
#### Install dependencies
```bash
go get
```

#### Run API server
```bash
go run main.go demo
```

#### Call query endpoint to return short-description (e.g)
```bash
curl --location --request GET 'localhost:8081/api?query=Yoshua_Bengio'
```
