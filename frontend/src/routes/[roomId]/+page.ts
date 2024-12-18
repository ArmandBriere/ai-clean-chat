/** @type {import('./$types').PageLoad} */
export function load({ params }: any) {
  return {
    roomId: params.roomId
  };
}
