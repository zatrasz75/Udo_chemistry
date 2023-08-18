// js.js
function calculateMolarMasses() {
    const nitrate = document.getElementById('nitrate').value;
    const phosphate = document.getElementById('phosphate').value;
    const potassium = document.getElementById('potassium').value;
    const micro = document.getElementById('micro').value;

    const data = {
        nitrate: nitrate,
        phosphate: phosphate,
        potassium: potassium,
        micro: micro
    };

    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/calculateMolarMasses", true);
    xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            const result = JSON.parse(xhr.responseText);
            document.getElementById("result").innerText = JSON.stringify(result, null, 2);
        }
    };

    xhr.send(JSON.stringify(data));
}
