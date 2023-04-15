import {get, readable} from "@square/svelte-store";

describe('given readable', async () => {
    it('should load once', async() => {

        let setCounter = 0
        let sut = readable<string | undefined>(undefined, (set) => {
            setCounter = setCounter + 1
            set('result')
        })

        console.log('setCounter', setCounter)
        console.log('beforeload', get(sut))
        console.log('setCounter', setCounter)
        console.log('loadresult', await sut.load())
        console.log('setCounter', setCounter)
        console.log('afterload', get(sut))
        console.log('setCounter2', setCounter)
        console.log('loadresult2', await sut.load())
        console.log('setCounter3', setCounter)

    })
})