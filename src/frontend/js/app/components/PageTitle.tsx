import React, { PropsWithChildren } from 'react';
import { Link } from 'react-router-dom';

type Props = PropsWithChildren<{
    action?: {
        title: string;
        disabled?: boolean;
    } & ( { href: string, onClick?: never } | { onClick: () => void, href?: never } )
}>

export const PageTitle = ({ action, children }: Props) => {
    return <div className="border-b border-gray-300 pb-5 sm:flex sm:items-center sm:justify-between">
        <h3 className="font-semibold leading-6 text-gray-900 text-xl">{ children }</h3>
        { action && <Action { ...action } /> }
    </div>
}


const Action = ({ disabled, title, href, onClick }: { disabled?: boolean, title: string, href?: string, onClick?: () => void}) => {
    const action = (() => {
        if (href) {
            return <Link
                to={ href }
                className="inline-flex items-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
            >
                { title }
            </Link>
        }

        if (onClick) {
            return <button
                disabled={ disabled }
                onClick={ onClick }
                className="inline-flex items-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:bg-gray-300 disabled:text-gray-500 disabled:shadow-none disabled:hover:bg-gray-300 disabled:hover:text-gray-500"
            >
                { title }
            </button>
        }

        return null
    })()

    if (!action) {
        return null
    }

    return <div className="mt-3 sm:ml-4 sm:mt-0">
        { action }
    </div>
}
