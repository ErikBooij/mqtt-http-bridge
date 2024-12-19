import React, { useContext, useEffect, useState } from 'react'
import { LayoutContext } from '../components/Layout';
import { Subscription } from '../../api/subscriptions';
import { Link } from 'react-router-dom';
import { PageTitle } from '../components/PageTitle';
import { useListSubscriptions } from '../../rq/subscription';

export const Subscriptions = () => {
    const { setCurrentPage } = useContext(LayoutContext);

    useEffect(() => {
        setCurrentPage('subscriptions');
    })

    const { isPending, error, data: subscriptions } = useListSubscriptions();

    if (isPending) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: { error.message }</div>;
    }

    if (!subscriptions?.length) {
        return <div>No subscriptions found</div>;
    }

    return (
        <>
            <PageTitle action={{ title: 'New Subscription', href: '/new-subscription' }}>Subscriptions</PageTitle>
            <div className="mt-8">
                <ul className="grid lg:grid-cols-2 gap-x-8">
                    { subscriptions.map(sub => (
                        <SubscriptionItem key={ sub.id } subscription={ sub } />
                    )) }
                </ul>
            </div>
        </>
    );
}

const SubscriptionItem = ({ subscription }: { subscription: Subscription }) => {
    const [ menuOpen, setMenuOpen ] = useState(false);

    return (
        <li className="flex items-center justify-between gap-x-6 py-5">
            <div className="min-w-0">
                <div className="flex items-start gap-x-3">
                    <p className="text-sm/6 font-semibold text-gray-900">{ subscription.name }</p>
                    <p className="mt-0.5 whitespace-nowrap rounded-md bg-green-50 px-1.5 py-0.5 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20">Active</p>
                </div>
                <div className="mt-1 flex items-center gap-x-2 text-xs/5 text-gray-500">
                    <p className="whitespace-nowrap js-subscription-id flex items-center gap-x-1 cursor-pointer">{ subscription.id }</p>
                    <svg viewBox="0 0 2 2" className="h-0.5 w-0.5 fill-current">
                        <circle cx="1" cy="1" r="1"/>
                    </svg>
                    <p className="truncate">{ subscription.topic }</p>
                </div>
            </div>
            <div className="flex flex-none items-center gap-x-4">
                <Link to={`/subscriptions/${subscription.id}`}
                   className="hidden rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:block">Edit</Link>
                <div className="relative flex-none">
                    <button type="button"
                            onClick={ () => setMenuOpen(!menuOpen) }
                            className="-m-2.5 block p-2.5 text-gray-500 hover:text-gray-900 js-subscription-menu-btn">
                        <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"
                             data-slot="icon">
                            <path
                                d="M10 3a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM10 8.5a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM11.5 15.5a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0Z"/>
                        </svg>
                    </button>
                    <div
                        className={ `absolute right-0 z-10 mt-2 w-32 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 focus:outline-none ${menuOpen ? '' : 'hidden'}` }>
                        <Link to={`/subscriptions/${subscription.id}`}
                           className="block px-3 py-1 text-sm/6 text-gray-900 sm:hidden" role="menuitem" tabIndex={ -1 }
                           id="options-menu-0-item-0">Edit</Link>
                        <Link to={ `/copy-subscription/${subscription.id}` }
                           className="block px-3 py-1 text-sm/6 text-gray-900" role="menuitem" tabIndex={ -1 }
                           id="options-menu-0-item-0">Duplicate</Link>
                        <Link to="#" className="block px-3 py-1 text-sm/6 text-red-900 js-delete-subscription"
                           role="menuitem" tabIndex={ -1 } id="options-menu-0-item-2">Delete</Link>
                    </div>
                </div>
            </div>
        </li>
    );
}
