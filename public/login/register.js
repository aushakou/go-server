window.addEventListener("load", () => {
    document.querySelector("button#register").addEventListener("click", () => {
        const username = document.querySelector("#username").value;
        const password = document.querySelector("#password").value;
        const baseUrl = `${window.location.protocol}//${window.location.host}`;
        const mode = document.querySelector("#mode").value;

        fetch(`${baseUrl}/api/register-${mode}`, {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({ username, password })
        })
            .then(res => {
                if (!res.ok) {
                    throw res.statusText;
                }
                window.location.assign(`${baseUrl}`);
            })
            .catch(err => {
                document.querySelector("#result").innerHTML = err;
                console.log(err);
            });
    });
});
