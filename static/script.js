async function createUUID() {
    var formData = new FormData();
    var extension = document.getElementById("upload-form").elements.namedItem("file").value.split(".").pop();
    formData.set("extension", extension)
    return await fetch(
       '/create?extension=' + extension,
       {
           method: "GET",
       })
}

async function uploadFile(uuid) {
    const file = document.querySelector('input[type="file"]');
    var formData = new FormData();
    formData.append('file', file.files[0]);
    return await fetch(
        '/upload?uuid=' + uuid,
        {
            method: "POST",
            body: formData
        }
    )
}

async function checkStatus(uuid) {
    return await fetch(
        '/check?uuid=' + uuid,
        {
            method: "GET"
        }
    )
}

function downloadFile(url, fileName){
    fetch(url, { method: 'get', mode: 'no-cors', referrerPolicy: 'no-referrer' })
        .then(res => res.blob())
        .then(res => {
            const aElement = document.createElement('a');
            aElement.setAttribute('download', fileName);
            const href = URL.createObjectURL(res);
            aElement.href = href;
            // aElement.setAttribute('href', href);
            aElement.setAttribute('target', '_blank');
            aElement.click();
            URL.revokeObjectURL(href);
        });
};

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function doSomething() {
    // Detect whether is file loaded
    if (document.getElementById("upload-form").elements.namedItem("file").value === "") {
        button.textContent = "Select a file to do stuff"
        return
    }

    button.textContent = "Doing stuff"
    const resUUID = await createUUID();
    const uuid = await resUUID.json();
    if (uuid.status === "ok") {
        const resUpload = await uploadFile(uuid.message);
        const ul = await resUpload.json();
        if (ul.status === "ok") {
            animation()
            while (true) {
                const resCheck = await checkStatus(uuid.message);
                const check = await resCheck.json();
                if (check.status === "wait")
                {
                    await sleep(500);
                } else {
                    // Get original file name
                    originalPathString = document.getElementById("upload-form").elements.namedItem("file").value.split(".").shift()
                    originalFilename = originalPathString.split(/(\\|\/)/g).pop()
                    downloadFile('/download?uuid=' + uuid.message, originalFilename + '.pdf')
                    button.textContent = "Done stuff"
                    await animationTerminator()
                    break
                }
            }
        } else {
            button.textContent = ul.message
            console.log("[E] Upload failed: " + ul.message )
        }
    } else {
        button.textContent = uuid.message
        console.log("[E] Error creating UUID: " + uuid.message)
    }
}

async function animationTerminator() {
    var button = document.getElementById('button')
    button.style.animation = "rotation 4s infinite linear, end 1s 1 ease-in-out"
}

async function animation() {
    var button = document.getElementById('button');
    button.style.animation = "rotation 4s infinite linear, goto 1s 1 ease-in-out"
    await sleep(1000)
    button.style.animation = "rotation 4s infinite linear, bouncing 4s infinite ease-in-out"
}