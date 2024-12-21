import React, { useContext, useEffect } from 'react'
import { LayoutContext } from '../components/Layout';
import { useListGlobalParameters } from '../../rq/parameter';
import { PageTitle } from '../components/PageTitle';
import { Link } from 'react-router-dom';

export const Parameters = () => {
    const { setCurrentPage } = useContext(LayoutContext);

    useEffect(() => {
        setCurrentPage('parameters');
    })

    const globalParameters = useListGlobalParameters();

    if (globalParameters.isPending) {
        return <div>Loading...</div>;
    }

    if (globalParameters.error) {
        return <div>Error: { globalParameters.error.message }</div>;
    }

    return (
        <div>
            <PageTitle
                action={ {
                    title: 'New Parameter',
                    href: '/new-parameter'
                } }
                currentPage="parameters"
            >
                Global Parameters
            </PageTitle>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-4 mt-4">
                { Object.entries(globalParameters.data).map(([ key, value ]) => (
                    <div
                        key={ key }
                        className="relative flex items-center space-x-3 rounded-lg border border-gray-300 bg-white px-6 py-5 shadow-sm focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 hover:border-gray-400"
                    >
                        <div className="min-w-0 flex-1">
                            <Link to={ '/parameters/' + key } className="focus:outline-none">
                                <p className="text-sm font-medium text-gray-900" title={ key }>{ key }</p>
                                <p className="truncate text-sm text-gray-500" title={ value }>{ value }</p>
                            </Link>
                        </div>
                    </div>
                )) }
            </div>
        </div>
    );
}
