#!/bin/bash

set -eo pipefail

file_path="${1}"
gopass_prefix="enpass"

declare -A directories

function translit() {
    local NAME=${*:-$(cat)};
    local TRS;
    TRS=$(sed "y/абвгдезийклмнопрстуфхцы/abvgdezijklmnoprstufxcy/" <<< "$NAME")
    TRS=$(sed "y/АБВГДЕЗИЙКЛМНОПРСТУФХЦЫ/ABVGDEZIJKLMNOPRSTUFXCY/" <<< "$TRS")
    TRS=${TRS//ч/ch};
    TRS=${TRS//Ч/CH} TRS=${TRS//ш/sh};
    TRS=${TRS//Ш/SH} TRS=${TRS//ё/jo};
    TRS=${TRS//Ё/JO} TRS=${TRS//ж/zh};
    TRS=${TRS//Ж/ZH} TRS=${TRS//щ/sh\'};
    TRS=${TRS///SH\'} TRS=${TRS//э/je};
    TRS=${TRS//Э/JE} TRS=${TRS//ю/ju};
    TRS=${TRS//Ю/JU} TRS=${TRS//я/ja};
    TRS=${TRS//Я/JA} TRS=${TRS//ъ/};
    TRS=${TRS//ъ} TRS=${TRS//ь/};
    TRS=${TRS//Ь/}
    hash iconv &> /dev/null && TRS=$(iconv -c -f UTF8 -t ASCII//TRANSLIT <<< "$TRS")
    echo ${TRS//[^[:alnum:]]/_} | awk '{print tolower($0)}';
}

function get_gopass_path() {
    local path_arr=()

    local category=$( echo ${1} | translit )
    local folder=$( echo ${2} | translit )
    local title=$( echo ${3} | translit )

    local trashed="${4:-0}"
    local archived="${5:-0}"
    local favorite="${6:-0}"

    path_arr+=( "$gopass_prefix" )
    if [[ "$trashed" == "1" ]]; then
        path_arr+=( "trash" )
    elif [[ "$archived" == "1" ]]; then
        path_arr+=( "archive" )
    elif [[ "$favorite" == "1" ]]; then
        path_arr+=( "favorite" )
    fi

    [ ! -z "$category" ] && path_arr+=( "${category}" )
    [ ! -z "$folder" ] && path_arr+=( "${folder}" )
    path_arr+=( "${title}" )

    echo $(IFS=/ ; echo "${path_arr[*]}")
}

function get_gopass_data_path() {
    local path=$(get_gopass_path ${@})
    echo "${path}/data";
}

function get_gopass_attachment_path() {
    local path=$(get_gopass_path ${@:1:$#-1})
    local attach_name=
    echo "${path}/attachments/${!#}";
}

function add_to_record() {
    local record="$1"
    local label=$( echo ${2} | translit )
    local value="${3}"
    local multiline="${4:-0}"

    if [ -z "${label}" -o -z "${value}" ]; then
        echo "${record}";
        return
    fi
    [ "$multiline" == "1" ] && {
        echo "${record}${label}: |\n  ${value}\n";
        return
    } || {
        echo "${record}${label}: ${value}\n";
        return
    }
    return
}

function jq_get_field() {
    echo $( echo ${1} | jq -r --arg selector "$2" '.[$selector]' );
}

while read -r folder; do
    uuid=$(echo ${folder} | jq -r '.uuid')
    title=$(echo ${folder} | jq -r '.title')
    directories[${uuid}]=$( echo ${title} | translit)
done <<< $(cat ${file_path} | jq -c '.folders // [] | .[] | {uuid, title}')


while read -r item; do
    record=""
    password=""

    archived=$(echo ${item} | jq -r '.archived')
    favorite=$(echo ${item} | jq -r '.favorite' )
    trashed=$(echo ${item} | jq -r '.trashed' )

    title=$(echo ${item} | jq -r '.title' )
    record=$(add_to_record "$record" "title" "${title}")

    subtitle=$(echo ${item} | jq -r '.subtitle' )
    record=$(add_to_record "$record" "subtitle" "${subtitle}")

    note=$(echo ${item} | jq -r '.note' )
    record=$(add_to_record "$record" "note" "${note}" "1")

    category=$(echo ${item} | jq -r '.category' )
    record=$(add_to_record "$record" "category" "${category}")

    template_type=$(echo ${item} | jq -r '.template_type' )
    record=$(add_to_record "$record" "template" "${template_type}")


    while read -r field; do
        deleted=$(echo ${field} | jq -r '.deleted')
        [[ "$deleted" == "1" ]] && continue

        value=$(echo ${field} | jq -r '.value')
        [ -z "$value" ] && continue

        type=$(echo ${field} | jq -r '.type')
        [[ "$type" == "section" ]] && continue

        label=$(echo ${field} | jq -r '.label' | translit )
        sensitive=$(echo ${field} | jq -r '.sensitive')

        multiline="0"
        [ "$type" == "totp" ] && label="totp"
        [ "$type" == "multiline" ] && multiline="1"

        if [ -z "$password" ] && [ "$type" == "password" ]; then
            password=${value}
        else
            record=$(add_to_record "$record" "${label}" "${value}" "${multiline}")
        fi

    done <<< $(echo ${item} | jq -c '.fields // [] | sort_by(.order) | .[]' )

    # writing password in record
    record="${password}\n---\n${record}"

    folders=()
    while read -r folder_uuid; do
        [ -z "${folder_uuid}" ] && continue
        folder_title=${directories[${folder_uuid}]}
        folders+=( "$folder_title" )
    done <<< $(echo ${item} | jq -r '.folders // [] | .[]')

    tags_val=$(IFS=, ; echo "${folders[*]}")
    record=$(add_to_record "$record" "tags" "[ ${tags_val} ]")

    while read -r attachment; do
        uuid=$(echo ${attachment} | jq -r '.uuid')
        [ -z "$uuid" ] && continue

        data=$(echo ${attachment} | jq -r '.data')
        kind=$(echo ${attachment} | jq -r '.kind')
        name=$(echo ${attachment} | jq -r '.name')
    done <<< $(echo ${item} | jq -c '.attachments // [] | sort_by(.order) | .[]')

    echo -e "$record"
    break
done <<< "$(cat ${file_path} | jq -c '.items // [] | .[]')"

# gopass insert --append
