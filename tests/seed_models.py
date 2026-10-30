"""
Seed script: creates one record of each model through the API (Feel-The-Night-Swagger)
Requirements: pip install requests

Usage:
    python3 seed_models.py
    BASE_URL=http://localhost:8080 python3 tests/seed_models.py

The sign-up endpoint does not return a token. The script logs in at /auth/login
and uses the access_token in the following requests. New accounts start as
"beginner"; content creation endpoints require "veteran" or "admin".
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
        print("!! Response is not JSON:", resp.text)
        sys.exit(1)
    if hide_response:
        print("Response received (content omitted).")
    else:
        print(json.dumps(data, indent=2, ensure_ascii=False))
    if not resp.ok:
        if hide_response:
            print(json.dumps(data, indent=2, ensure_ascii=False))
        print(f"!! Request failed with status {resp.status_code}")
        sys.exit(1)
    return data


def extract_value(data: dict, *keys):
    for k in keys:
        if k in data and data[k] is not None:
            return data[k]
    # also look inside "data" in case the response is wrapped
    if "data" in data and isinstance(data["data"], dict):
        for k in keys:
            if k in data["data"] and data["data"][k] is not None:
                return data["data"][k]
    return None


def main():
    print("==> Creating User...")
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
        print("!! Could not extract user_id from the response above. "
              "Adjust the keys in extract_value() and run again.")
        sys.exit(1)
    print(f"User created with ID: {user_id}\n")

    print("==> Logging in...")
    login_resp = post(LOGIN_PATH, {
        "email": credentials["email"],
        "password": credentials["password"],
    }, hide_response=True)
    access_token = extract_value(login_resp, "access_token")
    if access_token is None:
        print("!! Could not extract access_token from the login response.")
        sys.exit(1)
    SESSION.headers.update({"Authorization": f"Bearer {access_token}"})
    print("Logged in; token set for the following requests.\n")
    print("!! New users have the 'beginner' role; the API requires 'veteran' or "
          "'admin' to create characters, guides and events.\n")

    print("==> Creating Character...")
    char_resp = post(CHARACTERS_PATH, {
        "description": "test",
        "name": "test",
        "type": "shoto",
    })
    character_id = extract_value(char_resp, "CharacterID", "character_id", "id")
    if character_id is None:
        print("!! Could not extract character_id from the response above. "
              "Adjust the keys in extract_value() and run again.")
        sys.exit(1)
    print(f"Character created with ID: {character_id}\n")

    print("==> Creating Guide...")
    post(GUIDES_PATH, {
        "character_id": character_id,
        "description": "test",
        "link": "http://localhost:8080/swagger/index.html",
        "title": "test",
        "type": "combo",
        "user_id": user_id,
    })
    print()

    print("==> Creating Event...")
    post(EVENTS_PATH, {
        "day": "23/09/2026",
        "description": "test",
        "title": "test",
        "user_id": user_id,
    })
    print()

    print("==> Done! (LastEvent was not created because it is normally "
          "generated from an Event, not from a direct payload)")


if __name__ == "__main__":
    main()
