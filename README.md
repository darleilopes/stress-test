Gerando imagem:
```bash
docker build -t stresstest .
```

Executando o teste:
```bash
docker run stresstest --url=http://192.168.3.1/ --requests=10 --concurrency=10
```

Nota: 192.168.3.1 é meu roteador