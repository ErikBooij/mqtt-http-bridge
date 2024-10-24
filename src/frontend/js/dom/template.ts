export const elementFromTemplate = <T extends HTMLElement = HTMLElement>(template: string): T => {
    const tpl = document.createElement('template');
    tpl.innerHTML = template.trim();

    return tpl.content.firstChild as T;
}
