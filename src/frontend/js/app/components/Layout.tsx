import React, { createContext, useState } from 'react'
import { Link, Outlet } from 'react-router-dom';

const menuItems = [
    {
        id: 'subscriptions',
        title: 'Subscriptions',
        url: '/'
    },
    {
        id: 'parameters',
        title: 'Parameters',
        url: '/parameters'
    },
    {
        id: 'live-mqtt',
        title: 'Live MQTT',
        url: '/live-mqtt'
    }
] as const;

type LayoutContextData = {
    currentPage: typeof menuItems[number]['id'];
    setCurrentPage: (currentPage: typeof menuItems[number]['id']) => void;
};

export const LayoutContext = createContext<LayoutContextData>({
    currentPage: 'subscriptions',
    setCurrentPage: () => {}
});

export const Layout = () => {
    const [ currentPage, setCurrentPage ] = useState<LayoutContextData['currentPage']>('subscriptions');

    return (
        <LayoutContext.Provider value={{ currentPage, setCurrentPage }}>
            <nav className="bg-gray-800 fixed w-full top-0 left-0 right-0 px-4">
                <div className="mx-auto max-w-7xl">
                    <div className="relative flex h-16 items-center justify-between">
                        <div className="flex items-center justify-end sm:justify-start w-full">
                            <div className="hidden sm:block">
                                <div className="flex space-x-4">
                                    { menuItems.map((menuItem) => (
                                        <Link
                                            key={ menuItem.id }
                                            to={ menuItem.url }
                                            className={ `rounded-md px-3 py-2 text-sm font-medium ${ currentPage === menuItem.id ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white' }` }
                                        >
                                            { menuItem.title }
                                        </Link>
                                    )) }
                                </div>
                            </div>
                            <div className="-mr-2 flex sm:hidden">
                                <button type="button"
                                        className="relative inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-700 hover:text-white focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white js-menu-trigger"
                                        aria-controls="mobile-menu" aria-expanded="false">
                                    <span className="absolute -inset-0.5"></span>
                                    <span className="sr-only">Open main menu</span>
                                    <svg className="block h-6 w-6 js-menu-state-closed" fill="none" viewBox="0 0 24 24"
                                         strokeWidth="1.5" stroke="currentColor" aria-hidden="true">
                                        <path strokeLinecap="round" strokeLinejoin="round"
                                              d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"/>
                                    </svg>
                                    <svg className="hidden h-6 w-6 js-menu-state-open" fill="none" viewBox="0 0 24 24"
                                         strokeWidth="1.5" stroke="currentColor" aria-hidden="true">
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12"/>
                                    </svg>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="lg:hidden">
                    <div className="space-y-1 px-2 pb-3 pt-2 hidden js-menu">
                        { menuItems.map((menuItem) => (
                            <Link
                                key={ menuItem.id }
                                to={ menuItem.url }
                                className={ `block rounded-md px-3 py-2 text-base font-medium ${ currentPage === menuItem.id ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white' }` }
                            >
                                { menuItem.title }
                            </Link>
                        )) }
                    </div>
                </div>
            </nav>
            <div className="mx-auto max-w-7xl pt-20 px-4 xl:px-0">
                <Outlet/>
            </div>
        </LayoutContext.Provider>
    );
}
