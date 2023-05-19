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
    const resUUID = await createUUID();
    const uuid = await resUUID.json();
    if (uuid.status === "ok") {
        const resUpload = await uploadFile(uuid.message);
        const ul = await resUpload.json();
        if (ul.status === "ok") {
            while (true) {
                const resCheck = await checkStatus(uuid.message);
                const check = await resCheck.json();
                if (check.status === "wait")
                {
                    sleep(1000);
                } else {
                    // Get original file name
                    originalPathString = document.getElementById("upload-form").elements.namedItem("file").value.split(".").shift()
                    originalFilename = originalPathString.split(/(\\|\/)/g).pop()
                    downloadFile('/download?uuid=' + uuid.message, originalFilename + '.pdf')
                    break
                }
            }
        } else {
            console.log("[E] Upload failed: " + ul.message )
        }
    } else {
        console.log("[E] Error creating UUID: " + uuid.message)
    }
}