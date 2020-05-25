input = document.getElementById("Url")

async function solve(event) {
    event.preventDefault()

    // console.log(input.value) 
    await fetch("http://localhost:8000/a", {
        mode: "no-cors",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            value: input.value 
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.shorturl == "toang") {
            // console.log("???")
            document.getElementById("answer").innerHTML = "Your URL was wrong!!"
        } else {
            // console.log(data)
            document.getElementById("answer").innerHTML = "http://localhost:8000/s/" + data.shorturl
        }
        input.value = ""
    })
}