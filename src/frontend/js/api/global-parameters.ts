import {z} from "zod";
import {apiErrorFromResponse, AsyncMaybeAPIError, responseParseError} from "./common";
import {confirm, ConfirmStyle} from "../dom/confirm";
import {deleteSubscription} from "./subscriptions";

const globalParametersSchema = z.record(z.string(), z.string());
const globalParametersResponseSchema = z.object({
    parameters: globalParametersSchema,
}).strict();

export type GlobalParameters = z.infer<typeof globalParametersSchema>;

export const fetchGlobalParameters = async (): AsyncMaybeAPIError<GlobalParameters> => {
    const response = await fetch('/api/v1/global-parameters');

    if (response.status !== 200) {
        return [null, await apiErrorFromResponse(response)];
    }

    const parsedResponse = globalParametersResponseSchema.safeParse(await response.json());

    if (!parsedResponse.success) {
        return [null, responseParseError(parsedResponse.error)];
    }

    return [parsedResponse.data.parameters, null];
}

export const setGlobalParameter = async (key: string, value: string): AsyncMaybeAPIError<string> => {
    const response = await fetch('/api/v1/global-parameters', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({key, value}),
    });

    if (response.status !== 200) {
        return [null, await apiErrorFromResponse(response)];
    }

    return [key, null];
}

export const deleteGlobalParameter = async (key: string): AsyncMaybeAPIError<null> => {
    const response = await fetch(`/api/v1/global-parameters/${key}`, {
        method: 'DELETE',
    });

    if (response.status !== 200) {
        return [null, await apiErrorFromResponse(response)];
    }

    return [null, null];
}
