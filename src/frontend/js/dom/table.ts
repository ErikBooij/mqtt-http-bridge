import {elementFromTemplate} from "./template";

type TableEntry<T extends object = {}> = T

export interface Table<T extends TableEntry> {
    removeRow(select: (item: T) => boolean): void;
    removeRows(select: (item: T) => boolean): void;

    showLoading(): void;

    showError(error: string): void;

    showData(data: T[], columns?: TableRenderOptions<T>): void;
}

type TableRenderOptions<T extends TableEntry> = {
    actions?: (item: T) => TableActionItem<T>[]
    columns?: Columns<T>
    order: (keyof T)[],
}

type Columns<T extends TableEntry> = {
    [K in keyof T]?: Column<T>
}

type Column<T extends TableEntry> = {
    label?: string;
    mobile?: ColumnOptionMobile
    renderer?: (value: T) => string | HTMLElement
    width?: ColumnWidthOption
}

export enum ColumnOptionMobile {
    HIDE,
    COMBINE,
}

export enum ColumnWidthOption {
    MINIMAL,
}

type TableActionItem<T extends TableEntry> = {
    type?: 'destructive' | 'normal';
    label: string;
    action?: (item: T, table: Table<T>) => void;
    href?: string;
} & ({ action: (item: T, table: Table<T>) => void } | { href: string });

export const table = <T extends TableEntry>(selector: string, options?: TableRenderOptions<T>): Table<T> => {
    var t: Table<T>;

    let tableRows: { row: HTMLTableRowElement, item: T}[] = [];

    const element = document.querySelector(selector);

    if (!element) {
        throw new Error(`Element with selector ${selector} not found`);
    }

    element.classList.add('divide-y', 'divide-gray-300', 'table-fixed')

    const removeRows = (select: (item: T) => boolean, max?: number) => {
        let removed = 0;

        for (const { row, item } of tableRows) {
            if (select(item)) {
                row.remove();
                removed++;
            }

            if (max !== undefined && removed >= max) {
                break;
            }
        }
    }

    const removeRow = (select: (item: T) => boolean) => {
        removeRows(select, 1);
    }

    const showLoading = () => {
        element.innerHTML = 'Loading...';
    }

    const showError = (error: string) => {
        element.innerHTML = `Error: ${error}`;
    }

    const showData = (data: T[]) => {
        if (data.length === 0) {
            element.innerHTML = 'No data';
            return;
        }

        const cols = (options?.order ?? Object.keys(data[0])) as (keyof T)[]
        tableRows = [];

        element.innerHTML = '';
        element.append(
            createHeader(cols, options),
            createBody(cols, data, options, tableRows, t),
        );
    }

    return t = {
        removeRow,
        removeRows,
        showLoading,
        showError,
        showData,
    };
}

const createBody = <T extends TableEntry>(
    cols: (keyof T)[],
    data: T[],
    options: TableRenderOptions<T> | undefined,
    tableRows: { row: HTMLTableRowElement, item: T }[],
    t: Table<T>,
) => {
    const combinedColumns = cols.filter(c => options?.columns?.[c]?.mobile === ColumnOptionMobile.COMBINE)

    const body = elementFromTemplate(`<tbody class="divide-y divide-gray-200"></tbody>`)

    for (const row of data) {
        const entry = elementFromTemplate<HTMLTableRowElement>(`<tr></tr>`)

        let firstCol = true;

        for (const col of cols) {
            const cellContents = options?.columns?.[col]?.renderer?.(row) || (row[col] as string)
            const tableCell = elementFromTemplate(`<td></td>`)

            if (typeof cellContents === 'string') {
                tableCell.textContent = cellContents
            } else {
                tableCell.append(cellContents)
            }

            entry.append(tableCell)

            tableCell.classList.add(
                'py-4',
                'text-sm',
                'px-3',
                'text-gray-600',
                'truncate',
                'first:px-auto',
                'first:max-w-0',
                'first:py-4',
                'first:pl-4',
                'first:pr-3',
                'first:text-sm',
                'first:font-medium',
                'first:text-gray-900',
                'first:sm:max-w-none',
                'first:w-full',
                'first:sm:w-auto',
                'first:sm:pl-0',
            )

            if (!firstCol && isHiddenOnMobile(col, options)) {
                tableCell.classList.add('hidden', 'lg:table-cell')
            }

            if (options?.columns?.[col]?.width === ColumnWidthOption.MINIMAL) {
                tableCell.classList.add('w-0', 'whitespace-nowrap')
            }

            if (firstCol && combinedColumns.length) {
                const compactAttrs = elementFromTemplate(`<dl class="font-normal lg:hidden></dl>`)

                for (const combinedCol of combinedColumns) {
                    compactAttrs.append(elementFromTemplate(`<dt class="sr-only">${options?.columns?.[combinedCol]?.label ?? col as string}</dt>`))
                    compactAttrs.append(elementFromTemplate(`<dd class="mt-1 truncate text-gray-700">${row[combinedCol]}</dd>`))
                }

                tableCell.append(compactAttrs)
            }

            firstCol = false
        }

        if (hasActions(options)) {
            const actionCell = elementFromTemplate(`
                <td class="whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0 w-1">
                    <div class="flex justify-end gap-x-2 js-actions"></div>
                </td>
            `)

            for (const { label, type, href, action } of options?.actions?.(row) ?? []) {
                const actionButton = elementFromTemplate<HTMLAnchorElement>(`<a href="#" class="-my-2 p-2">${label}</a>`)

                if (type === 'destructive') {
                    actionButton.classList.add('text-red-600', 'hover:text-red-900')
                } else {
                    actionButton.classList.add('text-slate-600', 'hover:text-slate-900')
                }

                if (href) {
                    actionButton.setAttribute('href', href)
                } else if (action) {
                    actionButton.addEventListener('click', () => action(row, t))
                }

                actionCell.querySelector('.js-actions')!.append(actionButton)
            }

            entry.append(actionCell)
        }

        tableRows.push({ row: entry, item: row })

        body.append(entry)
    }

    return body
}

const createHeader = <T extends TableEntry>(cols: (keyof T)[], options?: TableRenderOptions<T>) => {
    const header = elementFromTemplate('<thead><tr></tr></thead>');
    const headerRow = header.querySelector('tr')!

    let firstCol = true;

    for (const col of cols) {
        const headerCell = elementFromTemplate(`<th scope="col">${options?.columns?.[col]?.label ?? col as string}</th>`)
        headerCell.classList.add(
            'py-3.5',
            'text-left',
            'text-sm',
            'font-semibold',
            'text-gray-900',
            'px-3',
            'first:px-auto',
            'first:pl-4',
            'first:pr-3',
            'first:sm:pl-0',
        )

        if (!firstCol && isHiddenOnMobile(col, options)) {
            headerCell.classList.add('hidden', 'lg:table-cell')
        }

        if (options?.columns?.[col]?.width === ColumnWidthOption.MINIMAL) {
            headerCell.classList.add('w-0', 'whitespace-nowrap')
        }

        headerRow.append(headerCell)

        firstCol = false
    }

    if (hasActions(options)) {
        headerRow.append(elementFromTemplate('<th scope="col"></th>'))
    }

    return header
}

const hasActions = <T extends TableEntry>(options?: TableRenderOptions<T>) => {
    return Boolean(options?.actions)
}

const isHiddenOnMobile = <T extends TableEntry>(col: keyof T, options: TableRenderOptions<T> | undefined) => {
    return options?.columns?.[col]?.mobile === ColumnOptionMobile.HIDE || options?.columns?.[col]?.mobile === ColumnOptionMobile.COMBINE
}
