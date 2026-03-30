export type Service = {
    container_id:string,
    created_at:string,
    description?:string,
    id:string,
    code?:string,
    image:string,
    language:string,
    name:string,
    status:string
}

export type ServiceCreateForm =  {
    code:string,
    description:string,
    language:string,
    name:string
}

export type ServiceUpdateForm = {
    code:string,
    description:string
}