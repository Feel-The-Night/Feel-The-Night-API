"""
Seed script: cria 1 registro de cada model via API (Feel-The-Night-Swagger)
Requisitos: pip install requests

Uso:
    python3 seed_models.py
    BASE_URL=http://localhost:8080 python3 tests/seed_models.py

O endpoint de criação não retorna token. O script faz login em /auth/login
e usa o access_token nas requisições seguintes. Contas novas começam como
"iniciante"; endpoints de criação de conteúdo exigem "veterano" ou "admin".
"""

import os
import sys
import json
import requests

BASE_URL = os.environ.get("BASE_URL", "http://localhost:8080")

USERS_PATH = "/users"
LOGIN_PATH = "/auth/login"
CHARACTERS_PATH = "/characters"
GUIDES_PATH = "/guides"
EVENTS_PATH = "/events"

SESSION = requests.Session()


def post(path: str, body: dict, *, hide_response: bool = False) -> dict:
    url = f"{BASE_URL}{path}"
    resp = SESSION.post(url, json=body)
    print(f"--> POST {url} [{resp.status_code}]")
    try:
        data = resp.json()
    except ValueError:
        print("!! Resposta não é JSON:", resp.text)
        sys.exit(1)
    if hide_response:
        print("Resposta recebida (conteúdo omitido).")
    else:
        print(json.dumps(data, indent=2, ensure_ascii=False))
    if not resp.ok:
        if hide_response:
            print(json.dumps(data, indent=2, ensure_ascii=False))
        print(f"!! Requisição falhou com status {resp.status_code}")
        sys.exit(1)
    return data


def extract_value(data: dict, *keys):
    for k in keys:
        if k in data and data[k] is not None:
            return data[k]
    # tenta dentro de "data" caso a resposta venha envelopada
    if "data" in data and isinstance(data["data"], dict):
        for k in keys:
            if k in data["data"] and data["data"][k] is not None:
                return data["data"][k]
    return None


def main():
    print("==> Criando User...")
    credentials = {
        "discord_id": "kisaltoo",
        "email": "kisalto@email.com",
        "nickname": "kisalto",
        "password": "Batata123",
    }
    user_resp = post(USERS_PATH, {
        "discord_id": credentials["discord_id"],
        "email": credentials["email"],
        "nickname": credentials["nickname"],
        "password": credentials["password"],
    })
    user_id = extract_value(user_resp, "UserID", "user_id", "id")
    if user_id is None:
        print("!! Não consegui extrair o user_id da resposta acima. "
              "Ajuste as chaves em extract_value() e rode de novo.")
        sys.exit(1)
    print(f"User criado com ID: {user_id}\n")

    print("==> Autenticando usuário...")
    login_resp = post(LOGIN_PATH, {
        "email": credentials["email"],
        "password": credentials["password"],
    }, hide_response=True)
    access_token = extract_value(login_resp, "access_token")
    if access_token is None:
        print("!! Não consegui extrair o access_token da resposta de login.")
        sys.exit(1)
    SESSION.headers.update({"Authorization": f"Bearer {access_token}"})
    print("Login concluído; token configurado para as próximas requisições.\n")
    print("!! Usuários novos têm perfil 'iniciante'; a API exige 'veterano' ou "
          "'admin' para criar personagens, guias e eventos.\n")

    print("==> Criando Character...")
    char_resp = post(CHARACTERS_PATH, {
        "description": "teste",
        "name": "teste",
        "type": "shoto",
    })
    character_id = extract_value(char_resp, "CharacterID", "character_id", "id")
    if character_id is None:
        print("!! Não consegui extrair o character_id da resposta acima. "
              "Ajuste as chaves em extract_value() e rode de novo.")
        sys.exit(1)
    print(f"Character criado com ID: {character_id}\n")

    print("==> Criando Guide...")
    post(GUIDES_PATH, {
        "character_id": character_id,
        "description": "teste",
        "link": "http://localhost:8080/swagger/index.html",
        "title": "teste",
        "type": "combo",
        "user_id": user_id,
    })
    print()

    print("==> Criando Event...")
    post(EVENTS_PATH, {
        "day": "23/09/2026",
        "description": "teste",
        "title": "teste",
        "user_id": user_id,
    })
    print()

    print("==> Concluído! (LastEvent não foi criado pois normalmente é "
          "gerado a partir de um Event, não via payload direto)")


if __name__ == "__main__":
    main()