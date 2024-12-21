import React, { PropsWithChildren, useContext, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { CurrentPage, LayoutContext } from './Layout';

type Props = PropsWithChildren<{
    action?: Action;
    secondaryAction?: Action;
    currentPage: CurrentPage;
}>

type Action = {
    title: string;
    disabled?: boolean;
} & ( { href: string, onClick?: never } | { onClick: () => void, href?: never } );

export const PageTitle = ({ action, children, currentPage, secondaryAction }: Props) => {
    const { setCurrentPage } = useContext(LayoutContext);

    useEffect(() => {
        setCurrentPage(currentPage);
    })

    return <div className="border-b border-gray-300 pb-5 sm:flex sm:items-center sm:justify-between">
        <h3 className="font-semibold leading-6 text-gray-900 text-xl">{ children }</h3>
        { (action || secondaryAction) && <div className="mt-3 sm:ml-4 sm:mt-0 flex gap-x-4">
            { secondaryAction && <Action { ...secondaryAction } type="secondary"/> }
            { action && <Action { ...action } type="primary"/> }
        </div> }
    </div>
}


const Action = ({ disabled, title, href, onClick, type }: Action & { type: 'primary' | 'secondary' }) => {
    const classes = type === 'primary'
        ? 'inline-flex items-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600'
        : 'inline-flex items-center rounded-md bg-slate-100 px-3 py-2 text-sm font-semibold text-slate-600 shadow-sm hover:bg-slate-300 hover:text-slate-800 border border-slate-300 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600';

    return ( () => {
        if (href) {
            return <Link to={ href } className={ classes }>{ title }</Link>
        }

        if (onClick) {
            return <button disabled={ disabled } onClick={ onClick } className={ classes }>{ title }</button>
        }

        return null
    } )()
}
