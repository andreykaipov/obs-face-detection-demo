#!/bin/sh

export TWITCH_CLIENT_ID=4lji6zvnpbkqn6c7nvz0nj9epm1qmo

user=povofkai
broadcaster_id="$(twitch api get users -q "login=$user" | jq -r .data[].id)"
current_tags="$(
        twitch api get /streams/tags -q "broadcaster_id=$broadcaster_id" |
                jq -r '.data[] | "\(.tag_id) \(.localization_names."en-us")"'
)"

# search for games:
# twitch api get /search/categories -Pu -q query='deep rock' | jq .

title="wow wow ow test"
game_id=494839

63e83904-a70b-4709-b963-a37a105d9932 Cooperative
1eba3cfe-51cc-460a-8259-bc8bb987f904 Competitive
cc8d5abb-39c9-4942-a1ee-e1558512119e Casual Playthrough

twitch api patch /channels -b '{
        "broadcaster_language": "en",
        "broadcaster_id": "'"$broadcaster_id"'",
        "game_id": "'"$game_id"'",
        "title": "'"$title"'"
}'

echo "$current_tags"
