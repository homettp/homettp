import type { Ref } from 'vue';

export interface BreadcrumbLink {
    name: string | Ref<string>
    href?: string
}

export interface CommandImage {
    name: string
    value: string
}

export interface FlashMessageItem {
    id: number
    value: string
    createdAt: number
}
