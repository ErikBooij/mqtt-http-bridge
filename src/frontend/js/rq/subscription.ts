import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
    createSubscription,
    fetchSubscription,
    fetchSubscriptions,
    Subscription,
    SubscriptionWithoutID,
    updateSubscription
} from '../api/subscriptions';
import { ApiRequestError, unpackMaybeAPIError } from './helpers';

const fetchSubscriptionQueryKey = (id: string) => [ 'fetchSubscription', id ];
const listSubscriptionsQueryKey = () => [ 'listSubscriptions' ];

export const useCreateSubscription = ({ onSuccess }: { onSuccess?: () => void}) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ subscription }: { subscription: SubscriptionWithoutID }): Promise<Subscription> => {
            return unpackMaybeAPIError(await createSubscription({ ...subscription }));
        },
        onError: (error: ApiRequestError) => {
            alert(error.errors.join('\n'));
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: listSubscriptionsQueryKey() });

            onSuccess?.();
        },
    })
}

export const useFetchSubscription = (id: string, enabled: boolean = true) => {
    return useQuery({
        queryKey: fetchSubscriptionQueryKey(id),
        queryFn: async () => {
            return unpackMaybeAPIError(await fetchSubscription(id));
        },
        enabled,
    })
}

export const useListSubscriptions = () => {
    return useQuery({
        queryKey: listSubscriptionsQueryKey(),
        queryFn: async () => {
            return unpackMaybeAPIError(await fetchSubscriptions());
        }
    })
}

export const useUpdateSubscription = ({ onSuccess }: { onSuccess?: () => void}) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, subscription }: { id: string, subscription: SubscriptionWithoutID }): Promise<Subscription> => {
            return unpackMaybeAPIError(await updateSubscription({ id, ...subscription }));
        },
        onSuccess: ({ id }) => {
            queryClient.invalidateQueries({ queryKey: listSubscriptionsQueryKey() });
            queryClient.invalidateQueries({ queryKey: fetchSubscriptionQueryKey(id) });

            onSuccess?.();
        },
    })
}
