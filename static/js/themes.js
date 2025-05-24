function setTheme(theme) {
    if (theme === 'pink') {
        document.documentElement.style.setProperty('--bg-gradient', 'linear-gradient(25deg, #ee628e 0%, #fde0e0 50%, #ec6590 100%)');
        document.documentElement.style.setProperty('--text-color', '#333333');
        document.documentElement.style.setProperty('--logo-bg', '#e77a9c');
        document.documentElement.style.setProperty('--logo-color', '#091833');
        document.documentElement.style.setProperty('--navbar-bg', 'rgba(0, 0, 0, 0.5)');
        document.documentElement.style.setProperty('--navbar-border', '#fff');
        document.documentElement.style.setProperty('--navlink-color', 'black');
        document.documentElement.style.setProperty('--navlink-hover', '#e77a9c');
    } else if (theme === 'aqua') {
        document.documentElement.style.setProperty('--bg-gradient', 'linear-gradient(25deg,rgb(63, 99, 156) 0%, #c2e9fb 50%,rgb(31, 100, 211) 100%)');
        document.documentElement.style.setProperty('--text-color', '#0f1b2a');
        document.documentElement.style.setProperty('--logo-bg', '#4facfe');
        document.documentElement.style.setProperty('--logo-color', '#ffffff');
        document.documentElement.style.setProperty('--navbar-bg', 'rgba(79, 172, 254, 0.4)');
        document.documentElement.style.setProperty('--navbar-border', '#4facfe');
        document.documentElement.style.setProperty('--navlink-color', '#0f1b2a');
        document.documentElement.style.setProperty('--navlink-hover', '#0077b6');
    } else if (theme === 'forest') {
        document.documentElement.style.setProperty('--bg-gradient', 'linear-gradient(25deg, #556270, #4ECDC4, #556270)');
        document.documentElement.style.setProperty('--text-color', '#f1f8e9');
        document.documentElement.style.setProperty('--logo-bg', '#2e7d32');
        document.documentElement.style.setProperty('--logo-color', '#ffffff');
        document.documentElement.style.setProperty('--navbar-bg', 'rgba(46, 125, 50, 0.4)');
        document.documentElement.style.setProperty('--navbar-border', '#81c784');
        document.documentElement.style.setProperty('--navlink-color', '#f1f8e9');
        document.documentElement.style.setProperty('--navlink-hover', '#aed581');
    } else if (theme === 'dark') {
        document.documentElement.style.setProperty('--bg-gradient', 'linear-gradient(25deg, #2f2f2f,rgb(116, 132, 139), #4f4f4f)');
        document.documentElement.style.setProperty('--text-color', '#f0f0f0');
        document.documentElement.style.setProperty('--logo-bg', '#3f3f3f');
        document.documentElement.style.setProperty('--logo-color', '#e0e0e0');
        document.documentElement.style.setProperty('--navbar-bg', 'rgba(50, 50, 50, 0.7)');
        document.documentElement.style.setProperty('--navbar-border', '#6f6f6f');
        document.documentElement.style.setProperty('--navlink-color', '#c0c0c0');
        document.documentElement.style.setProperty('--navlink-hover', '#e8e8e8');
    }
}
