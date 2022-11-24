import parser from 'uri-template'
interface Expandable {
    expand(values: Record<string, unknown>): string
}
export const expand = (template: string, params: {}): string => {
    let tpl = parser.parse(template)
    return tpl.expand(params)
}
export const expandTpl = (parsed: Expandable, params: {}): string => {
    return parsed.expand(params)
}