import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { deleteGlobalParameter, fetchGlobalParameters, setGlobalParameter, } from '../api/global-parameters';
import { ApiRequestError, unpackMaybeAPIError } from './helpers';

const listGlobalParametersQueryKey = () => [ 'listGlobalParameters' ];

export const useSetGlobalParameter = ({ onSuccess }: { onSuccess?: () => void }) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ key, value }: { key: string, value: string }): Promise<string> => {
            return unpackMaybeAPIError(await setGlobalParameter(key, value));
        },
        onError: (error: ApiRequestError) => {
            alert(error.errors.join('\n'));
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: listGlobalParametersQueryKey() });

            onSuccess?.();
        },
    })
}

export const useListGlobalParameters = () => {
    return useQuery({
        queryKey: listGlobalParametersQueryKey(),
        queryFn: async () => {
            return unpackMaybeAPIError(await fetchGlobalParameters());
        },
    })
}

export const useDeleteGlobalParameter = ({ onSuccess }: { onSuccess?: () => void }) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ key }: { key: string }): Promise<null> => {
            return unpackMaybeAPIError(await deleteGlobalParameter(key));
        },
        onError: (error: ApiRequestError) => {
            alert(error.errors.join('\n'));
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: listGlobalParametersQueryKey() });

            onSuccess?.();
        },
    })
}
