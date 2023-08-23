// js.js
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
        if (xhr.readyState === 4 && xhr.status === 200) {
            const result = JSON.parse(xhr.responseText);
            document.getElementById("result").innerText = JSON.stringify(result, null, 2);
        }
    };

    xhr.send(JSON.stringify(data));
}