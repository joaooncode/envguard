---
sidebar_position: 2
title: Instalação
---

# Instalação

O `envguard` é distribuído como um binário compilado único em Go, leve e sem dependências externas de runtime.

---

## 1. Via `go install` (Recomendado para desenvolvedores Go)

Se você já possui o **Go 1.22+** instalado no seu ambiente:

```bash
go install github.com/joaooncode/envguard/cmd/envguard@latest
```

Certifique-se de que o diretório `$GOPATH/bin` (ou `$HOME/go/bin` no Linux/macOS, `%USERPROFILE%\go\bin` no Windows) esteja configurado na sua variável de ambiente `PATH`.

---

## 2. Compilando a partir do Código-Fonte

Você pode clonar o repositório e compilar diretamente:

```bash
# Clone o repositório
git clone https://github.com/joaooncode/envguard.git
cd envguard

# Compile o binário
go build -o envguard ./cmd/envguard
```

Para mover para um diretório no PATH do sistema:

### Linux / macOS

```bash
sudo mv envguard /usr/local/bin/
```

### Windows (PowerShell como Administrador)

```powershell
Move-Item .\envguard.exe C:\Windows\System32\
```

---

## 3. Verificando a Instalação

Após a instalação, verifique se a CLI está acessível executando:

```bash
envguard version
```

Saída esperada:

```text
envguard v0.1.0
```

---

## 4. Distribuições Futuras (No Roadmap)

Estamos trabalhando para disponibilizar gerenciadores de pacotes nativos:

- **Homebrew** (`brew install joaooncode/tap/envguard`)
- **Scoop** (`scoop install envguard`)
- **WinGet** (`winget install joaooncode.envguard`)
- **GitHub Releases** (binários pré-compilados para Linux, macOS e Windows x86_64 / arm64)
