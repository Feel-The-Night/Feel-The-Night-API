"""
Seed script: cria 1 registro de cada model via API (Feel-The-Night-Swagger)
Requisitos: pip install requests

Uso:
    python3 seed_models.py
    BASE_URL=http://localhost:8080 python3 "tests/Seed models.py"
"""

import os
import sys
import json
import requests

BASE_URL = os.environ.get("BASE_URL", "http://localhost:8080")

USERS_PATH = "/users"
CHARACTERS_PATH = "/characters"
GUIDES_PATH = "/guides"
EVENTS_PATH = "/events"


def post(path: str, body: dict) -> dict:
    url = f"{BASE_URL}{path}"
    resp = requests.post(url, json=body)
    print(f"--> POST {url} [{resp.status_code}]")
    try:
        data = resp.json()
    except ValueError:
        print("!! Resposta não é JSON:", resp.text)
        sys.exit(1)
    print(json.dumps(data, indent=2, ensure_ascii=False))
    if not resp.ok:
        print(f"!! Requisição falhou com status {resp.status_code}")
        sys.exit(1)
    return data


def extract_id(data: dict, *keys):
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
    user_resp = post(USERS_PATH, {
        "discord_id": "kisaltoo",
        "email": "kisalto@email.com",
        "nickname": "kisalto",
        "password": "Batata123",
    })
    user_id = extract_id(user_resp, "UserID", "user_id", "id")
    if user_id is None:
        print("!! Não consegui extrair o user_id da resposta acima. "
              "Ajuste as chaves em extract_id() e rode de novo.")
        sys.exit(1)
    print(f"User criado com ID: {user_id}\n")

    print("==> Criando Character...")
    char_resp = post(CHARACTERS_PATH, {
        "description": "teste",
        "name": "teste",
        "type": "shoto",
    })
    character_id = extract_id(char_resp, "CharacterID", "character_id", "id")
    if character_id is None:
        print("!! Não consegui extrair o character_id da resposta acima. "
              "Ajuste as chaves em extract_id() e rode de novo.")
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