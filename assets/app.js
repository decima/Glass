import "./app.css";


let theme = localStorage.getItem('theme');
try {
    document.querySelector("html").classList.remove(theme)

    if (theme === "dark") {
        document.querySelector("html").classList.add(theme)

    }
} catch (error) {
    //do nothing : we know there was expecting theme from localstorage but nothing.
}

document.querySelector(".theme-toggler").addEventListener("click", () => {
    document.querySelector("html").classList.remove(theme)
    if (theme === "dark") {
        theme = "light"
    } else {
        theme = "dark"
    }
    document.querySelector("html").classList.add(theme)
    localStorage.setItem('theme', theme);

})
