function calculateMolarMasses() {
    const nitrateMass = document.getElementById('nitrate_mass').value;
    const phosphateMass = document.getElementById('phosphate_mass').value;
    const potassiumMass = document.getElementById('potassium_mass').value;
    const microMass = document.getElementById('micro_mass').value;

    const nitrate = document.getElementById('nitrate').value;
    const phosphate = document.getElementById('phosphate').value;
    const potassium = document.getElementById('potassium').value;
    const micro = document.getElementById('micro').value;

    const data = {
        nitrate: nitrate,
        phosphate: phosphate,
        potassium: potassium,
        micro: micro,
        nitrate_mass: nitrateMass,
        phosphate_mass: phosphateMass,
        potassium_mass: potassiumMass,
        micro_mass: microMass
    };

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/", true);
    xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                // Обновляем содержимое <div id="result"> с данными из ответа сервера
                document.getElementById("result").innerHTML = xhr.responseText;
            } else {
                // Обработка ошибок, если необходимо
                console.error("Ошибка при выполнении запроса:", xhr.status);
            }
        }
    }

    xhr.send(JSON.stringify(data));
}

function deleteRecord() {
    const recordId = document.getElementById('recordId').value;

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/delet?id=" + recordId, true);

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                // Обновляем содержимое <div id="result"> с сообщением об успешном удалении
                document.getElementById("result").innerHTML = "Запись успешно удалена.";
            } else {
                // Обработка ошибок, если необходимо
                console.error("Ошибка при выполнении запроса:", xhr.status);
            }
        }
    }

    xhr.send();
}

