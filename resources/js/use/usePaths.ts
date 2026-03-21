export default () => {
    const commandIndexPath = () => '/';

    const commandCreatePath = () => '/command/create';

    const commandEditPath = (id: number) => `/command/edit?id=${id}`;

    const commandRefreshTokenPath = (id: number) => `/command/refresh-token?id=${id}`;

    const commandDeletePath = (id: number) => `/command/delete?id=${id}`;

    const callHistoryPath = (commandID: number) => `/call/history?commandID=${commandID}`;

    const callDeletePath = (id: number) => `/call/delete?id=${id}`;

    const userIndexPath = () => '/user';

    const userCreatePath = () => '/user/create';

    const userEditPath = (id: number) => `/user/edit?id=${id}`;

    const userDeletePath = (id: number) => `/user/delete?id=${id}`;

    const loginPath = () => '/login';

    const logoutPath = () => '/logout';

    return {
        callDeletePath,
        callHistoryPath,
        commandCreatePath,
        commandDeletePath,
        commandEditPath,
        commandIndexPath,
        commandRefreshTokenPath,
        loginPath,
        logoutPath,
        userCreatePath,
        userDeletePath,
        userEditPath,
        userIndexPath
    };
};
