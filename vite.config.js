export default {
    server: {
        watch: {
            disableGlobbing: false,
        },
        fs: {
            strict: false,
            allow: [".."],
        },
    },
    root: "./assets",
    base: "/build/",
    build: {
        manifest: true,
        emptyOutDir: true,
        assetsDir: "",
        outDir: "../public/build/",
        rollupOptions: {
            input: "./assets/app.js",
        },
    }
}