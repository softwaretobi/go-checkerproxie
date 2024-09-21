Proxy Checker

Un outil simple pour vérifier la validité des proxies depuis un fichier texte.

![image](https://github.com/user-attachments/assets/a5477516-75d5-4d32-b830-54c81648d017)

Fonctionnalités

    Vérification rapide des proxies
    Enregistrement des proxies valides dans un fichier de sortie
    Supporte les proxies HTTP/HTTPS

Utilisation 
- `go run proxy_checker.go <proxy_file> <proxy_type> <output_file> <timeout>`


Exemple : 
- `go run proxy_checker.go proxies.txt http valides.txt 5000`
