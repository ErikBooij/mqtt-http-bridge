const selectors = {
    menu: '.js-menu',
    menuTrigger: '.js-menu-trigger',
    menuStateOpen: '.js-menu-state-open',
    menuStateClosed: '.js-menu-state-closed',
} as const;

export const initializeMenu = () => {
    let menuOpen = false;

    const menu = document.querySelector(selectors.menu);
    const menuTrigger = document.querySelector(selectors.menuTrigger);
    const menuStateOpen = document.querySelector(selectors.menuStateOpen);
    const menuStateClosed = document.querySelector(selectors.menuStateClosed);

    if (!menu || !menuTrigger || !menuStateOpen || !menuStateClosed) {
        console.groupCollapsed(`Not all required menu elements found.`)
        console.log(`menu: ${menu}`)
        console.log(`menuTrigger: ${menuTrigger}`)
        console.log(`menuStateOpen: ${menuStateOpen}`)
        console.log(`menuStateClosed: ${menuStateClosed}`)
        console.groupEnd()

        return
    }

    menuTrigger.addEventListener('click', () => {
        menuOpen = !menuOpen;

        if (menuOpen) {
            menu.classList.remove('hidden');
            menuStateOpen.classList.remove('hidden');
            menuStateClosed.classList.add('hidden');
        } else {
            menu.classList.add('hidden');
            menuStateOpen.classList.add('hidden');
            menuStateClosed.classList.remove('hidden');
        }
    })
}
