import type { PostItType } from '~/types'

const POSTIT_TYPES: PostItType[] = [
  {
    kind: 'event',
    label: 'Evento de Domínio',
    description: 'Algo relevante que aconteceu no domínio (verbo no passado).',
    color: '#F59E0B',
    textColor: '#1F2937'
  },
  {
    kind: 'command',
    label: 'Comando',
    description: 'Ação ou intenção disparada por um ator ou política.',
    color: '#3B82F6',
    textColor: '#FFFFFF'
  },
  {
    kind: 'actor',
    label: 'Ator / Usuário',
    description: 'Pessoa ou papel que dispara comandos.',
    color: '#FACC15',
    textColor: '#1F2937'
  },
  {
    kind: 'aggregate',
    label: 'Agregado',
    description: 'Cluster de regras de negócio responsável por comandos e eventos.',
    color: '#FDE68A',
    textColor: '#1F2937'
  },
  {
    kind: 'policy',
    label: 'Política',
    description: 'Reação automática: quando X acontece, então Y.',
    color: '#A78BFA',
    textColor: '#1F2937'
  },
  {
    kind: 'readmodel',
    label: 'Read Model / View',
    description: 'Informação consultada por atores para decidir.',
    color: '#34D399',
    textColor: '#1F2937'
  },
  {
    kind: 'external',
    label: 'Sistema Externo',
    description: 'Sistema fora do domínio que produz ou consome eventos.',
    color: '#F472B6',
    textColor: '#1F2937'
  },
  {
    kind: 'hotspot',
    label: 'Hot Spot',
    description: 'Dúvida, conflito ou ponto a esclarecer.',
    color: '#EF4444',
    textColor: '#FFFFFF'
  }
]

const POSTIT_TYPE_MAP: Record<string, PostItType> = POSTIT_TYPES.reduce((acc, type) => {
  acc[type.kind] = type
  return acc
}, {} as Record<string, PostItType>)

export const usePostItTypes = () => {
  const types = POSTIT_TYPES
  const byKind = (kind: string): PostItType => POSTIT_TYPE_MAP[kind] ?? POSTIT_TYPES[0]
  return { types, byKind }
}
